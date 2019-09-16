from .tool.func import *

def user_tool_2(name):
    
    
    data = '''
        <h2>''' + load_lang('tool') + '''</h2>
        <ul>
            <li><a href="/record/''' + url_pas(name) + '''">''' + load_lang('record') + '''</a></li>
            <li><a href="/topic/user:''' + url_pas(name) + '''">''' + load_lang('user_discussion') + '''</a></li>
        </ul>
    '''
            
    if admin_check(1) == 1:
        sqlQuery("select block from ban where block = ?", [name])
        if sqlQuery("fetchall"):
            ban_name = load_lang('ban_release')
        else:
            ban_name = load_lang('ban')
    
        data += '''
            <h2>''' + load_lang('admin') + '''</h2>
            <ul>
                <li><a href="/ban/''' + url_pas(name) + '''">''' + ban_name + '''</a></li>
                <li><a href="/check/''' + url_pas(name) + '''">''' + load_lang('check') + '''</a></li>
            </ul>
        '''

    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('tool') + ')', 0])],
        data = data,
        menu = 0
    ))