from .tool.func import *

def give_history_add_2(conn, name):
    curs = conn.cursor()

    ip = ip_check()
    if admin_check() != 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        today = get_time()
        content = flask.request.form.get('content', '')
        content = savemark(content)
        leng = '+' + str(len(content))

        history_plus(
            name,
            content,
            today,
            'Add:' + flask.request.form.get('get_ip', ''),
            flask.request.form.get('send', ''),
            leng
        )

        conn.commit()

        return redirect('/history/' + url_pas(name))
    else:
        curs.execute(db_change('select data from other where name = "edit_bottom_text"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            b_text = '<hr class=\"main_hr\">' + sql_d[0][0]
        else:
            b_text = ''

        curs.execute(db_change('select data from other where name = "edit_help"'))
        sql_d = curs.fetchall()
        if sql_d and sql_d[0][0] != '':
            p_text = sql_d[0][0]
        else:
            p_text = load_lang('defalut_edit_help')

        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('history_add'), wiki_set(), custom(), other2([' (' + name + ')', 0])],
            data = '''
                <form method="post">
                    <script>do_stop_exit();</script>
                    ''' + edit_button() + '''
                    <textarea rows="25" id="content" placeholder="''' + p_text + '''" name="content"></textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('name') + '''" name="get_ip" type="text">
                    <hr class=\"main_hr\">
                    <button id="save" type="submit" onclick="go_save_zone = 1;">''' + load_lang('save') + '''</button>
                    <button id="preview" type="button" onclick="load_preview(\'''' + url_pas(name) + '\')">' + load_lang('preview') + '''</button>
                </form>
                ''' + b_text + '''
                <hr class=\"main_hr\">
                <div id="see_preview"></div>
            ''',
            menu = [['history/' + url_pas(name), load_lang('return')]]
        ))