#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import datetime

import json
import random
import time

import websockets


@dataclass
class pixel:
    x: int
    y: int
    color: int
    timestamp: int
    userid: int

async def sender():
    async with websockets.connect("ws://localhost:8080/") as websocket:
        while True:
            message = pixel(
                x=random.randint(0, 9),
                y=random.randint(0, 9),
                color=random.randint(0,15),
                timestamp=int(time.time()),
                userid=1,
            )
            await websocket.send(json.dumps(message.__dict__))
            await asyncio.sleep(0.1)

async def client():
    async with websockets.connect("ws://localhost:8080/") as websocket:
        i= 0
        while True:
            i+=1
            x = await websocket.recv()
            print(i, pixel(**json.loads(x)))

async def main():
    coros = [sender() for _ in range(100)]
    coros.append(client())
    returns = await asyncio.gather(*coros)

if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
