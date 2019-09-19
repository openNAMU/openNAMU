from .tool.func import *

def recent_history_tool_2(name):
    

    num = str(int(number_check(flask.request.args.get('num', '1'))))

    data = '''
        <h2>''' + load_lang('tool') + '''</h2>
        <ul>
            <li>
                <a href="/raw/''' + url_pas(name) + '?num=' + num + '">' + load_lang('raw') + '''</a> 
            </li>
    '''

    if (int(num) - 1) > 0:
        data += '''
            <li>
                <a href="/diff/''' + url_pas(name) + '?first=''' + str(int(num) - 1) + '&second=' + num + '">' + load_lang('compare') + '''</a>
            </li>
        '''

    if flask.request.args.get('type', '') == 'history':
        data += '''
            <li>
                <a href="/revert/''' + url_pas(name) + '?num=' + num + '">' + load_lang('revert') + '''</a>
            </li>
        '''
    elif (int(num) - 1) > 0:
        data += '''
            <li>
                <a href="/revert/''' + url_pas(name) + '?num=' + str(int(num) - 1) + '">' + load_lang('revert') + '''</a>
            </li>
        '''

    if admin_check(6) == 1:
        sqlQuery("select title from history where title = ? and id = ? and hide = 'O'", [name, num])
        hide = sqlQuery("fetchall")
        data += '''
            <li>
                <a href="/hidden/''' + url_pas(name) + '?num=' + num + '">' + (load_lang('hide_release') if hide else load_lang('hide')) + '''
            </li>
        '''

    if admin_check() == 1:
        data += '''
            <li>
                <a href="/history_delete/''' + url_pas(name) + '?num=' + num + '">' + load_lang('history_delete') + '''
            </li>
        '''

    data += '''
        </ul>
    '''
                
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2(['(r' + num + ')', 0])],
        data = data,
        menu = [['history/' + url_pas(name), load_lang('return')]]
    ))