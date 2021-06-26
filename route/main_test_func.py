import time
from .tool.func import *

def main_test_func_2(conn):
    test_start = time.time()

    for _ in range(0, 10000):
        load_lang('edit')
        
    end_time = "time :" + str(time.time() - test_start) + '\n'

    return end_time