from .tool.func import *

def main_error_404(e = ''):
    if os.path.exists('404.html') and flask.request.path != '/':
        return open('404.html', encoding = 'utf8').read(), 404
    else:
        return redirect('/w/' + wiki_set(2))