from .tool.func import *

def login_find():
    with get_db_connect() as conn:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('password_search'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <ul class="opennamu_ul">
                    <li><a href="/login/find/email">''' + load_lang('email') + '''</a></li>
                    <li><a href="/login/find/key">''' + load_lang('key') + '''</a></li>
                </ul>
            ''',
            menu = [['user', load_lang('return')]]
        ))