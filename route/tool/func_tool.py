import urllib.parse
import datetime
import hashlib
import flask
import re

import os
import html
import json
import sqlite3
import threading
import ipaddress

set_data = ''

def get_time():
    return str(datetime.datetime.today().strftime("%Y-%m-%d %H:%M:%S"))

def db_data_get(data):
    global set_data
    
    set_data = data

def db_change(data):
    if set_data == 'mysql':
        data = data.replace('random()', 'rand()')
        data = data.replace('%', '%%')
        data = data.replace('?', '%s')

    return data

def url_pas(data):
    return urllib.parse.quote(data).replace('/','%2F')

def sha224_replace(data):
    return hashlib.sha224(bytes(data, 'utf-8')).hexdigest()

def md5_replace(data):
    return hashlib.md5(data.encode()).hexdigest()

def ip_in_range(ip, ip_range):
    return ipaddress.ip_address(ip) in ipaddress.ip_network(ip_range)