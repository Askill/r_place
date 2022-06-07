#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import datetime
from PIL import Image
import json
from multiprocessing import Pool
import random
import time
import matplotlib
from matplotlib import pyplot as plt
import numpy as np

import websockets
import cv2
import matplotlib.image as mpimg


@dataclass
class pixel:
    x: int
    y: int
    color: int
    timestamp: int
    userid: int

def hex_to_rgb(h):
    return tuple(int(h[i:i+2], 16) for i in (0, 2, 4))

hex_colors = [
        "#FFFFFF",
        "#E4E4E4",
        "#888888",
        "#222222",
        "#FFA7D1",
        "#E50000",
        "#E59500",
        "#A06A42",
        "#E5D900",
        "#94E044",
        "#02BE01",
        "#00D3DD",
        "#0083C7",
        "#0000EA",
        "#CF6EE4",
        "#820080"
        ]
rgb_colors = [hex_to_rgb(h[1:]) for h in hex_colors]

def eucleadian_distance(rgb1, rgb2):
    if len(rgb1) != len(rgb2):
        raise ValueError
    sum_part = np.sum([(rgb1[i]-rgb2[i])**2 for i in range(len(rgb1))])
    # return np.sqrt(sum_part) # technically correct, but we only care about rank not exact distance and sqrt is expensive
    return sum_part

def closest_match(rgb, color_map):
    return min(range(len(rgb_colors)), key=lambda i: eucleadian_distance(rgb, color_map[i]))

async def sender(img):
    async with websockets.connect("ws://localhost:8080/set") as websocket:
        while True:
            rx = random.randint(0, 999)
            ry = random.randint(0, 999)
            message = pixel(
                x=rx,
                y=ry,
                color= closest_match(img[rx][ry], rgb_colors),
                timestamp=int(time.time()),
                userid=1,
            )
            await websocket.send(json.dumps(message.__dict__))
            succ = await websocket.recv()
            if succ == "1":
                print(message, "was not set")
            
            

async def client():
    image = np.zeros(shape=[1000, 1000, 3], dtype=np.uint8)
    async with websockets.connect("ws://localhost:8080/get") as websocket:
        i= 0
        while True:
            i+=1
            x = pixel(**json.loads(await websocket.recv()))
            image[x.x][x.y] = rgb_colors[x.color]
            #if i% 5000 == 0:
            #    cv2.imshow("changes x", image)
            #    cv2.waitKey(10) & 0XFF
            await websocket.send("1")
            #print(i, x)

async def main():
    img= Image.open('./2.jpg')
    img= img.resize((1000, 1000), Image.ANTIALIAS)
    img = np.array(img)
    coros = [sender(img) for _ in range(100)]
    coros.append(client())
    returns = await asyncio.gather(*coros)
    
def asyncMain(x):
    asyncio.get_event_loop().run_until_complete(main())

if __name__ == "__main__":
    with Pool(12) as p:
        print(p.map(asyncMain, [() for _ in range(12)]))
    #asyncMain(0)
