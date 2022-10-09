from .tool.func import *

def main_func_easter_egg():
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(skin_check(),
            imp = ['Easter Egg', wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '<iframe width="640" height="360" src="https://www.youtube.com/embed/iQPbgD_CTd4" frameborder="0" allowfullscreen></iframe>',
            menu = [['manager', load_lang('return')]]
        ))