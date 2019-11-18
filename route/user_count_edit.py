from .tool.func import *

def user_count_edit_2(conn, name):
    curs = conn.cursor()

    if name == None:
        that = ip_check()
    else:
        that = name

    curs.execute(db_change("select count(title) from history where ip = ?"), [that])
    count = curs.fetchall()
    if count:
        data = count[0][0]
    else:
        data = 0

    curs.execute(db_change("select count(title) from topic where ip = ?"), [that])
    count = curs.fetchall()
    if count:
        t_data = count[0][0]
    else:
        t_data = 0

    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('count'), wiki_set(), custom(), other2([0, 0])],
        data =  '''
                <ul>
                    <li><a href="/record/''' + url_pas(that) + '''">''' + load_lang('edit_record') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/topic_record/''' + url_pas(that) + '''">''' + load_lang('discussion_record') + '''</a> : ''' + str(t_data) + '''</a></li>
                </ul>
                ''',
        menu = [['user', load_lang('return')]]
    ))