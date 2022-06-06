from .tool.func import *

def main_func_setting_robot():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/ban')
        
        if flask.request.method == 'POST':
            curs.execute(db_change("select name from other where name = 'robot'"))
            if curs.fetchall():
                curs.execute(db_change("update other set data = ? where name = 'robot'"), [flask.request.form.get('content', '')])
            else:
                curs.execute(db_change("insert into other (name, data) values ('robot', ?)"), [flask.request.form.get('content', '')])

            conn.commit()

            fw = open('./robots.txt', 'w', encoding='utf8')
            fw.write(re.sub('\r\n', '\n', flask.request.form.get('content', '')))
            fw.close()

            admin_check(None, 'edit_set (robot)')

            return redirect('/setting/robot')
        else:
            if not os.path.exists('robots.txt'):
                curs.execute(db_change('select data from other where name = "robot"'))
                robot_test = curs.fetchall()
                if robot_test:
                    fw_test = open('./robots.txt', 'w', encoding='utf8')
                    fw_test.write(re.sub('\r\n', '\n', robot_test[0][0]))
                    fw_test.close()
                else:
                    fw_test = open('./robots.txt', 'w', encoding='utf8')
                    fw_test.write('User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/')
                    fw_test.close()

                    curs.execute(db_change('insert into other (name, data) values ("robot", "User-agent: *\nDisallow: /\nAllow: /$\nAllow: /w/")'))

            curs.execute(db_change("select data from other where name = 'robot'"))
            robot = curs.fetchall()
            if robot:
                data = robot[0][0]
            else:
                data = ''

            f = open('./robots.txt', encoding='utf8')
            lines = f.readlines()
            f.close()

            if not data or data == '':
                data = ''.join(lines)

            return easy_minify(flask.render_template(skin_check(),
                imp = ['robots.txt', wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = '''
                    <a href="/robots.txt">(''' + load_lang('view') + ''')</a>
                    <hr class="main_hr">
                    <form method="post">
                        <textarea rows="25" name="content">''' + html.escape(data) + '''</textarea>
                        <hr class="main_hr">
                        <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    </form>
                ''',
                menu = [['setting', load_lang('return')]]
            ))