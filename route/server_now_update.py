from .tool.func import *

def server_now_update_2(conn, r_ver):
    curs = conn.cursor()

    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'update')

        curs.execute('select data from other where name = "update"')
        up_data = curs.fetchall()
        if up_data:
            up_data = up_data[0][0]
        else:
            up_data = 'stable'

        print('----')
        print('Update')
        print('----')
        if platform.system() == 'Linux':
            ok = []

            ok += [os.system('git remote rm origin')]
            ok += [os.system('git remote add origin https://github.com/2DU/opennamu.git')]
            ok += [os.system('git fetch origin ' + up_data)]
            ok += [os.system('git reset --hard origin/' + up_data)]
            if (ok[0] and ok[1] and ok[2] and ok[3]) == 0:
                return redirect('/restart')
            else:
                print('Update error')
        elif platform.system() == 'Windows':
            os.system('rd /s /q route')
            urllib.request.urlretrieve('https://github.com/2DU/opennamu/archive/' + up_data + '.zip', 'update.zip')
            zipfile.ZipFile('update.zip').extractall('')
            ok = os.system('xcopy /y /s /r opennamu-' + up_data + ' .')
            if ok == 0:
                os.system('rd /s /q opennamu-' + up_data)
                os.system('del update.zip')

                return redirect('/restart')
            else:
                print('Update error')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('update'), wiki_set(), custom(), other2([0, 0])],
            data = load_lang("update_error") + ' <a href="https://github.com/2DU/opennamu">(Github)</a>',
            menu = [['manager/1', load_lang('return')]]
        ))
    else:
        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('update'), wiki_set(), custom(), other2([0, 0])],
            data = load_lang('update_warring') + '''
                <hr class=\"main_hr\">
                <ul>
                    <li>''' + load_lang('version') + ' : ' + r_ver + '''</li>
                    <li id="ver_send" style="display: none;">''' + load_lang('lastest') + ''' : </li>
                </ul>
                <a href="https://github.com/2du/openNAMU">(Master)</a> <a href="https://github.com/2du/openNAMU/tree/stable">(Stable)</a>
                <hr class=\"main_hr\">
                <form method="post">
                    <button type="submit">''' + load_lang('update') + '''</button>
                </form>
                <script>load_ver();</script>
            ''',
            menu = [['manager', load_lang('return')]]
        ))

