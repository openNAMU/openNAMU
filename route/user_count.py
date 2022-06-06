from .tool.func import *

def user_count(name = None):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if name == None:
            that = ip_check()
        else:
            that = name

        curs.execute(db_change("select count(*) from history where ip = ?"), [that])
        count = curs.fetchall()
        if count:
            data = count[0][0]
        else:
            data = 0

        curs.execute(db_change("select count(*) from topic where ip = ?"), [that])
        count = curs.fetchall()
        if count:
            t_data = count[0][0]
        else:
            t_data = 0

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('count'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
            data = '''
                <ul class="inside_ul">
                    <li><a href="/record/''' + url_pas(that) + '''">''' + load_lang('edit_record') + '''</a> : ''' + str(data) + '''</li>
                    <li><a href="/record/topic/''' + url_pas(that) + '''">''' + load_lang('discussion_record') + '''</a> : ''' + str(t_data) + '''</a></li>
                </ul>
            ''',
            menu = [['user', load_lang('return')]]
        ))