import zipfile
import urllib.request

from .tool.func import *

def main_sys_update():
    with get_db_connect() as conn:
        curs = conn.cursor()

        if admin_check(conn) != 1:
            return re_error(conn, '/error/3')

        if flask.request.method == 'POST':
            admin_check(conn, None, 'update')

            curs.execute(db_change('select data from other where name = "update"'))
            up_data = curs.fetchall()
            up_data = up_data[0][0] if up_data and up_data[0][0] in ['stable', 'beta', 'dev'] else 'stable'

            print('Update')
            
            if platform.system() == 'Linux':
                ok = []
                ok += [os.system('git remote rm origin')]
                ok += [os.system('git remote add origin https://github.com/opennamu/opennamu.git')]
                ok += [os.system('git fetch --depth=1 origin ' + up_data)]
                ok += [os.system('git reset --hard origin/' + up_data)]
                for for_a in ok[1:]:
                    if for_a != 0:
                        break
                else:
                    return redirect(conn, '/restart')
                
                print('Error : update failed')
            elif platform.system() == 'Windows':
                os.system('rd /s /q route')

                urllib.request.urlretrieve('https://github.com/opennamu/opennamu/archive/' + up_data + '.zip', 'update.zip')
                zipfile.ZipFile('update.zip').extractall('')
                ok = os.system('xcopy /y /s /r opennamu-' + up_data + ' .')
                if ok == 0:
                    os.system('rd /s /q opennamu-' + up_data)
                    os.system('del update.zip')

                    return redirect(conn, '/restart')
                else:
                    print('Error : update failed')

            return re_error(conn, '/error/34')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'update'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                data = get_lang(conn, 'update_warning') + '''
                    <hr class="main_hr">
                    <ul class="opennamu_ul">
                        <li id="ver_send_2">''' + get_lang(conn, 'version') + ''' : </li>
                        <li id="ver_send">''' + get_lang(conn, 'lastest') + ''' : </li>
                    </ul>
                    <a href="https://github.com/openNAMU/openNAMU">(Beta)</a> <a href="https://github.com/openNAMU/openNAMU/tree/stable">(Stable)</a>
                    <hr class="main_hr">
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'update') + '''</button>
                    </form>
                    <!-- JS : opennamu_do_insert_version -->
                ''',
                menu = [['manager', get_lang(conn, 'return')]]
            ))

