#!/usr/bin/env python

import asyncio
from dataclasses import dataclass
import os
from socket import timeout
from PIL import Image
import json
from multiprocessing import Pool
import random
import time
from matplotlib import pyplot as plt
import numpy as np

import websockets

@dataclass
class pixel:
    x: int
    y: int
    color: int
    timestamp: int
    userid: int


def hex_to_rgb(h):
    return tuple(int(h[i : i + 2], 16) for i in (0, 2, 4))


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
    "#820080",
]
rgb_colors = [hex_to_rgb(h[1:]) for h in hex_colors]


def eucleadian_distance(rgb1, rgb2):
    if len(rgb1) != len(rgb2):
        raise ValueError
    sum_part = np.sum([(i-j) ** 2 for i, j in zip(rgb1, rgb2)])
    # return np.sqrt(sum_part) # technically correct, but we only care about rank not exact distance and sqrt is expensive
    return sum_part


def closest_match(rgb, color_map):
    return min(
        range(len(rgb_colors)), key=lambda i: eucleadian_distance(rgb, color_map[i])
    )


async def sender(target, img):
    start_x = random.randint(0, 900)
    start_y = random.randint(0, 900)
    max_w, max_h, _ = img.shape
    i = 0
    async for websocket in websockets.connect(target + "/set", timeout=10):
        try:
            while i < max_h*max_w*1.8:
                i+=1
                rx = random.randint(0, max_w - 1)
                ry = random.randint(0, max_h - 1)
                if rx + start_x >= 1000 or ry + start_y >= 1000:
                    continue
                message = pixel(
                    x=rx + start_x,
                    y=ry + start_y,
                    color=closest_match(img[rx][ry], rgb_colors),
                    timestamp=int(time.time()),
                    userid=1,
                )
                await websocket.send(json.dumps(message.__dict__))
                succ = await websocket.recv()
                if succ != "0":
                    print(message, "was not set")
            return 
        except websockets.ConnectionClosed:
            print("reconnecting")
            continue


async def client(target):
    image = np.zeros(shape=[1000, 1000, 3], dtype=np.uint8)
    async for websocket in websockets.connect(target + "/get"):
        try:
            x = pixel(**json.loads(await websocket.recv()))
            image[x.x][x.y] = rgb_colors[x.color]
            await websocket.send("1")
        except websockets.ConnectionClosed:
            continue


def rescale(max_dimension, img):
    w, h = img.size
    maxi = max([w, h])
    scale = max_dimension / maxi
    return img.resize((int(scale * w), int(scale * h)), Image.ANTIALIAS)


async def main(target):
    images_folder_path = "./images"
    while True:
        print("start cycle")
        image_paths = os.listdir(images_folder_path)
        images = [
            np.array(rescale(random.randint(100, 400), Image.open(f"{images_folder_path}/{image_path}")))
            for image_path in image_paths
        ]
        coros = [sender(target, images[i % len(images)]) for i in range(len(images))]
        _ = await asyncio.gather(*coros)


def asyncMain(x, target):
    asyncio.get_event_loop().run_until_complete(main(target))


if __name__ == "__main__":
    # with Pool(12) as p:
    #    print(p.map(asyncMain, [() for _ in range(12)]))
    asyncMain(0, target="ws://" + os.getenv("TARGET", "venus:8080"))

