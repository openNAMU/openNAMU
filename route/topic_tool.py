from .tool.func import *

def topic_tool_2(name, sub):
    

    data = ''

    sqlQuery("select stop, agree from rd where title = ? and sub = ?", [name, sub])
    close_data = sqlQuery("fetchall")
    if close_data:
        if close_data[0][0] == 'S':
            t_state = 'Stop'
        elif close_data[0][0] == 'O':
            t_state = 'Close'
        else:
            t_state = 'Normal'
    else:
        t_state = 'Normal'

    if admin_check(3) == 1:
        data = '''
            <h2>''' + load_lang('admin_tool') + '''</h2>
            <ul>
                <li><a href="/topic/''' + url_pas(name) + '/sub/' + url_pas(sub) + '''/setting"> Topic setting</a></li>
            </ul>
        '''

    print(close_data)
    data += '''
        <h2>''' + load_lang('tool') + '''</h2>
        <ul>
            <li><a id="reload" href="javascript:void(0);" onclick="req_alarm();">''' + load_lang('use_push_alarm') + '''</a></li>
            <li>''' + load_lang('topic_state') + ''' : ''' + t_state + '' + (' (Agree)' if close_data and (close_data[0][1] == 'O') else '') + '''</li>
        </ul>
    '''

    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('topic_tool') + ')', 0])],
        data = data,
        menu = [['topic/' + url_pas(name) + '/sub/' + url_pas(sub), load_lang('return')]]
    ))