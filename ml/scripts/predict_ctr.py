#!/usr/bin/env python3
__author__ = "nicoroble"
__version__ = "0.1.0"
__license__ = "MIT"

import sys
import getopt
import random

def usage():
    print('-t = bool\tused to run multiple threads and store results if true')
    print('-r = bool\tused to read resulting csv files')
    print('-n = int\tspecify number of threads for each optimization algorithm')
    print('-f = string\tspecify which csv file to read')

def str2bool(v):
    return v.lower() in ("true")

def main(opts):
    """ Main entry point of the app """

    arguments = {}
    floor_price= int(),
    media_id= str(),
    timestamp= int(),
    os_type= str(),
    banner_size= int(),
    banner_position= int(),
    device_type= int(),
    gender= str(),
    age= int(),
    income= int(),
    has_child= str()
    is_married= str()

    for opt, arg in opts:
        if opt in ('-h', '--help'):
            usage()
            sys.exit(2)
        elif opt in ('-fp', '--floorPrice'):
            floor_price = int(arg)
            # print(floor_price)
        elif opt in ('-mi', '--mediaId'):
            media_id = arg
            # print(media_id)
        elif opt in ('-ts', '--timestamp'):
            timestamp = int(arg)
            # print(timestamp)
        elif opt in ('-ot', '--osType'):
            os_type = arg
            # print(os_type)
        elif opt in ('-bs', '--bannerSize'):
            banner_size = int(arg)
            # print(banner_size)
        elif opt in ('-bp', '--bannerPosition'):
            banner_position = int(arg)
            # print(banner_position)
        elif opt in ('-dt', '--deviceType'):
            device_type = arg
            # print(device_type)
        elif opt in ('-gd', '--gender'):
            gender = arg
            # print(gender)
        elif opt in ('-ag', '--age'):
            age = int(arg)
            # print(age)
        elif opt in ('-ic', '--income'):
            income = int(arg)
            # print(income)
        elif opt in ('-hc', '--hasChild'):
            has_child = arg
            # print(has_child)
        elif opt in ('-im', '--isMarried'):
            is_married = arg
            # print(is_married)
        else:
            usage()
            sys.exit(2)

    result = random.uniform(0, 1)
    return 3
if __name__ == '__main__':
    try:
        opts, args = getopt.getopt(
            sys.argv[1:],
            'fp:mi:ts:ot:bs:bp:dt:gd:ag:ic:hc:im:h',
            ['floorPrice=',
             'mediaId=',
             'timestamp=',
             'osType=',
             'bannerSize=',
             'bannerPosition',
             'deviceType',
             'gender=',
             'age=',
             'income=',
             'hasChild=',
             'isMarried=',
             'help'
             ]
        )
    except getopt.GetoptError:
        usage()
        sys.exit(2)

    main(opts)