from .tool.func import *

def application_submitted_2(conn):
    curs = conn.cursor()

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('application_submitted'), wiki_set(), custom(), other2([0, 0])],
        data =  '''<p>''' + load_lang('waiting_for_approval') + '''</p>''',
        menu = [['user', load_lang('return')]]
    ))