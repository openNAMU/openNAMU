from .tool.func import *

def search_2():
    

    return redirect('/search/' + url_pas(flask.request.form.get('search', 'test')))
