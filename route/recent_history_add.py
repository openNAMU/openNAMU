from .tool.func import *

from .edit import edit_editor

def recent_history_add(name = 'Test', do_type = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if admin_check() != 1:
            return re_error('/ban')

        if flask.request.method == 'POST':
            admin_check(None, 'history_add (' + name + ')')

            today = get_time()
            content = flask.request.form.get('content', '')
            leng = '+' + str(len(content))

            history_plus(
                name,
                content,
                today,
                'Add:' + flask.request.form.get('get_ip', ''),
                flask.request.form.get('send', ''),
                leng,
                t_check = 'add',
                mode = 'add'
            )

            conn.commit()

            return redirect('/history/' + url_pas(name))
        else:            
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('history_add'), wiki_set(), wiki_custom(), wiki_css(['(' + name + ')', 0])],
                data = '''
                    <form method="post">
                        <input placeholder="''' + load_lang('why') + '''" name="send">
                        <hr class="main_hr">
                        
                        <input placeholder="''' + load_lang('name') + '''" name="get_ip">
                        <hr class="main_hr">

                        ''' + edit_editor(curs, ip) + '''
                    </form>
                ''',
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))