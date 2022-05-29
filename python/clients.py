#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import datetime

import json
import random
import time
import matplotlib
from matplotlib import pyplot as plt
import numpy as np

import websockets
import cv2

@dataclass
class pixel:
    x: int
    y: int
    color: int
    timestamp: int
    userid: int

async def sender():
    async with websockets.connect("ws://localhost:8080/set") as websocket:
        while True:
            message = pixel(
                x=random.randint(0, 99),
                y=random.randint(0, 99),
                color=random.randint(0,11),
                timestamp=int(time.time()),
                userid=1,
            )
            await websocket.send(json.dumps(message.__dict__))
            
            await asyncio.sleep(0.1)

async def client():
    image = np.zeros(shape=[100, 100, 3], dtype=np.uint8)
    colors = []
    for name, hex in matplotlib.colors.cnames.items():
        colors.append(matplotlib.colors.to_rgb(hex))

    async with websockets.connect("ws://localhost:8080/get") as websocket:
        i= 0
        while True:
            i+=1
            x = pixel(**json.loads(await websocket.recv()))
            image[x.x][x.y] = ([y*255 for y in colors[x.color]])
            if i% 1000 == 0:
                cv2.imshow("changes x", image)
                cv2.waitKey(10) & 0XFF
            print(i, x)

async def main():
    coros = [sender() for _ in range(100)]
    coros.append(client())
    returns = await asyncio.gather(*coros)
    

if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
