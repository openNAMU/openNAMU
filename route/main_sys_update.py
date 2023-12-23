import zipfile
import urllib.request

from .tool.func import *

def main_sys_update():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check() != 1:
            return re_error('/error/3')

        if flask.request.method == 'POST':
            admin_check(None, 'update')

            curs.execute(db_change('select data from other where name = "update"'))
            up_data = curs.fetchall()
            up_data = up_data[0][0] if up_data and up_data[0][0] in ['stable', 'beta', 'dev'] else 'stable'

            print('Update')
            
            if platform.system() == 'Linux':
                ok = []

                ok += [os.system('git remote rm origin')]
                ok += [os.system('git remote add origin https://github.com/opennamu/opennamu.git')]
                ok += [os.system('git fetch origin ' + up_data)]
                ok += [os.system('git reset --hard origin/' + up_data)]
                if (ok[0] and ok[1] and ok[2] and ok[3]) == 0:
                    return redirect('/restart')
                else:
                    print('Error : update failed')
            elif platform.system() == 'Windows':
                os.system('rd /s /q route')
                urllib.request.urlretrieve('https://github.com/opennamu/opennamu/archive/' + up_data + '.zip', 'update.zip')
                zipfile.ZipFile('update.zip').extractall('')
                ok = os.system('xcopy /y /s /r opennamu-' + up_data + ' .')
                if ok == 0:
                    os.system('rd /s /q opennamu-' + up_data)
                    os.system('del update.zip')

                    return redirect('/restart')
                else:
                    print('Error : update failed')

            return re_error('/error/34')
        else:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('update'), wiki_set(), wiki_custom(), wiki_css([0, 0])],
                data = load_lang('update_warning') + '''
                    <hr class="main_hr">
                    <ul class="opennamu_ul">
                        <li id="ver_send_2">''' + load_lang('version') + ''' : </li>
                        <li id="ver_send">''' + load_lang('lastest') + ''' : </li>
                    </ul>
                    <a href="https://github.com/openNAMU/openNAMU">(Beta)</a> <a href="https://github.com/openNAMU/openNAMU/tree/stable">(Stable)</a>
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + load_lang('update') + '''</button>
                    </form>
                    <!-- JS : opennamu_do_insert_version -->
                ''',
                menu = [['manager', load_lang('return')]]
            ))

