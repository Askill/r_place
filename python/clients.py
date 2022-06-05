#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import datetime

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

async def sender(img):
    async with websockets.connect("ws://localhost:8080/set") as websocket:
        while True:
            rx = random.randint(0, 999)
            ry = random.randint(0, 999)
            message = pixel(
                x=rx,
                y=ry,
                color=int(sum(img[rx][ry])/3),
                timestamp=int(time.time()),
                userid=1,
            )
            await websocket.send(json.dumps(message.__dict__))
            succ = await websocket.recv()
            if succ == "1":
                print(message, "was not set")
            
            

async def client():
    image = np.zeros(shape=[1000, 1000, 3], dtype=np.uint8)
    colors = []
    for name, hex in matplotlib.colors.cnames.items():
        colors.append(matplotlib.colors.to_rgb(hex))

    async with websockets.connect("ws://localhost:8080/get") as websocket:
        i= 0
        while True:
            i+=1
            x = pixel(**json.loads(await websocket.recv()))
            #image[x.x][x.y] = ([y*255 for y in colors[x.color]])
            image[x.x][x.y] = ((x.color, x.color, x.color))
            #if i% 5000 == 0:
            #    cv2.imshow("changes x", image)
            #    cv2.waitKey(10) & 0XFF
            await websocket.send("1")
            #print(i, x)

async def main():
    img=mpimg.imread('./1.jpg')
    coros = [sender(img) for _ in range(100)]
    coros.append(client())
    returns = await asyncio.gather(*coros)
    
def asyncMain(x):
    asyncio.get_event_loop().run_until_complete(main())

if __name__ == "__main__":
    with Pool(12) as p:
        print(p.map(asyncMain, [() for _ in range(12)]))
    #asyncMain(0)
