#!/usr/bin/env python3

__author__ = "nicoroble"
__version__ = "0.1.0"
__license__ = "MIT"
import os
def main():
    """ Main entry point of the app """
    result = os.system("python -c predict_ctr.py --timestamp 0")
    print(result)

if __name__ == "__main__":
    """ This is executed when run from the command line """
    main()