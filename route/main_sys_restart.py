from .tool.func import *

def main_sys_restart_do():
    print('Restart')

    time.sleep(3)

    python_ver = ''
    python_ver = str(sys.version_info.major) + '.' + str(sys.version_info.minor)

    run_list = [
        sys.executable,
        'python' + python_ver,
        'python3',
        'python',
        'py -' + python_ver
    ]

    for exe_name in run_list:
        try:
            subprocess.Popen([exe_name] + sys.argv)
            break
        except:
            continue
    
    os._exit(0)

async def main_sys_restart(golang_process):
    with get_db_connect() as conn:
        if await acl_check('', 'owner_auth', '', '') == 1:
            return await re_error(conn, 3)

        if flask.request.method == 'POST':
            await acl_check(tool = 'owner_auth', memo = 'restart')

            if golang_process.poll() is None:
                golang_process.terminate()
                try:
                    golang_process.wait(timeout = 5)
                except subprocess.TimeoutExpired:
                    golang_process.kill()
                    try:
                        golang_process.wait(timeout = 5)
                    except subprocess.TimeoutExpired:
                        print('Golang process not terminated properly.')

            threading.Thread(target = main_sys_restart_do).start()
            return flask.Response(get_lang(conn, "warning_restart"), status = 200)
        else:
            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [get_lang(conn, 'wiki_restart'), await wiki_set(), await wiki_custom(conn), wiki_css([0, 0])],
                data = '''
                    <form method="post">
                        <button type="submit">''' + get_lang(conn, 'restart') + '''</button>
                    </form>
                ''',
                menu = [['manager', get_lang(conn, 'return')]]
            ))