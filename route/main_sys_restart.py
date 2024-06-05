from .tool.func import *

def main_sys_restart():
    with get_db_connect() as conn:
        if admin_check(conn) != 1:
            return re_error(conn, '/error/3')

        if flask.request.method == 'POST':
            admin_check(conn, None, 'restart')

            print('Restart')

            python_ver = ''
            python_ver = str(sys.version_info.major) + '.' + str(sys.version_info.minor)

            run_list = [sys.executable, 'python' + python_ver, 'python3', 'python', 'py -' + python_ver]
            for exe_name in run_list:
                try:
                    os.execl(exe_name, sys.executable, *sys.argv)
                except:
                    pass

                try:
                    os.execl(exe_name, '"' + sys.executable + '"', *sys.argv)
                except:
                    pass

                try:
                    os.execl(exe_name, os.path.abspath(__file__), *sys.argv)
                except:
                    pass
            else:
                return re_error(conn, '/error/33')
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'wiki_restart'), wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <div id="opennamu_main_sys_restart"></div>
                    <script defer src="/views/main_css/js/route/main_sys_restart.js''' + cache_v() + '''"></script>
                    <script>window.addEventListener("DOMContentLoaded", function() { opennamu_main_sys_restart(); });</script>
                ''',
                menu = [['manager', get_lang(conn, 'return')]]
            ))