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


async def main():

    async with websockets.connect("ws://localhost:8080/") as websocket:
        for i in range(10):
            message = pixel(
                x=random.randint(0, 10),
                y=random.randint(0, 10),
                color=random.randint(0,15),
                timestamp=int(time.time()),
                userid=1,
            )
            print(message)
            await websocket.send(json.dumps(message.__dict__))
        print(await websocket.recv())
        while True:
            x = await websocket.recv()
            print(x)


if __name__ == "__main__":
    asyncio.get_event_loop().run_until_complete(main())
