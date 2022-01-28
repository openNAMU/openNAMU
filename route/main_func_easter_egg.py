from .tool.func import *

def main_func_easter_egg(conn):
    return easy_minify(flask.render_template(skin_check(),
        imp = ['easter_egg.html', wiki_set(), wiki_custom(), wiki_css([0, 0])],
        data = open('./views/main_css/file/easter_egg.html', encoding='utf8').read(),
        menu = 0
    ))