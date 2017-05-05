from flask import Flask, request, session, render_template, send_file
app = Flask(__name__)

from urllib import parse
import json
import pymysql
import time
import re
import bcrypt
import os
import difflib
import hashlib

from func import *
from mark import *

json_data = open('set.json').read()
data = json.loads(json_data)