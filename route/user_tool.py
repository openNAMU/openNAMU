from .tool.func import *

def user_tool_2(conn, name):
    curs = conn.cursor()

    data = '''
        <h2>''' + load_lang('tool') + '''</h2>
        <ul>
            <li><a href="/record/''' + url_pas(name) + '''">''' + load_lang('record') + '''</a></li>
            <li><a href="/topic/user:''' + url_pas(name) + '''">''' + load_lang('user_discussion') + '''</a></li>
        </ul>
    '''

    if admin_check(1) == 1:
        curs.execute(db_change("select block from rb where block = ? and ongoing = '1'"), [name])
        ban_name = load_lang('release') if curs.fetchall() else load_lang('ban')
        
        data += '''
            <h2>''' + load_lang('admin') + '''</h2>
            <ul>
                <li><a href="/ban/''' + url_pas(name) + '''">''' + ban_name + '''</a></li>
                <li><a href="/check/''' + url_pas(name) + '''">''' + load_lang('check') + '''</a></li>
            </ul>
        '''

    return easy_minify(flask.render_template(skin_check(),
        imp = [name, wiki_set(), custom(), other2(['(' + load_lang('tool') + ')', 0])],
        data = data,
        menu = 0
    ))