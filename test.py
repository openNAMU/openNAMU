from flask import Flask, render_template
app = Flask(__name__)

import json
import pymysql

json_data=open('set.json').read()
data = json.loads(json_data)

@app.route('/')
def redirect():
    return '<meta http-equiv="refresh" content="0;url=/w/" />'

@app.route('/w/')
def w(name=None):
    return render_template('index.html')

if __name__ == '__main__':
    app.run()
