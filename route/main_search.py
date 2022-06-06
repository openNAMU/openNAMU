from .tool.func import *

def main_search():
    return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))
