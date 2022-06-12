from .tool.func import *

def recent_history_add(name = 'Test'):
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
                        <script>do_stop_exit();</script>
                        ''' + edit_button() + '''
                        <textarea rows="25" id="content" name="content"></textarea>
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                        <hr class="main_hr">
                        <input placeholder="''' + load_lang('name') + '''" name="get_ip" type="text">
                        <hr class="main_hr">
                        <button id="save" type="submit" onclick="go_save_zone = 1;">''' + load_lang('save') + '''</button>
                        <button id="preview" type="button" onclick="load_preview(\'''' + url_pas(name) + '\')">' + load_lang('preview') + '''</button>
                    </form>
                    <hr class="main_hr">
                    <div id="see_preview"></div>
                ''',
                menu = [['history/' + url_pas(name), load_lang('return')]]
            ))