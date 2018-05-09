# -*- coding: utf-8 -*-
#import fileinput

if __name__ == '__main__':
    print("start consume data")
    data = ""
    try:
        while data != "EOF":
            data = input()
            print(data, "|")
    except EOFError:
        print("cache EOFError")
    # for line in fileinput.input():
    #     if line == "EOF":
    #         break
    #     print(line, "|")
    print("consume complete")
