from .tool.func import *

def main_func_easter_egg():
    with get_db_connect() as conn:
        curs = conn.cursor()

        ip = ip_check()
        if ip_or_user(ip) == 0:
            curs.execute(db_change('select name from user_set where id = ? and name = ?'), [ip, 'get_🥚'])
            if not curs.fetchall():
                curs.execute(db_change('insert into user_set (name, id, data) values ("get_🥚", ?, "Y")'), [ip])
    
        if platform.system() == 'Linux':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.bin"), sys._getframe().f_code.co_name], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.bin"), sys._getframe().f_code.co_name], stdout = subprocess.PIPE).communicate()[0]

            data = data.decode('utf8')
        elif platform.system() == 'Windows':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.exe"), sys._getframe().f_code.co_name], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.exe"), sys._getframe().f_code.co_name], stdout = subprocess.PIPE).communicate()[0]

            data = data.decode('utf8')
        else:
            data = ''

        return easy_minify(conn, flask.render_template(skin_check(conn),
            imp = ['Easter Egg', wiki_set(conn), wiki_custom(conn), wiki_css([0, 0])],
            data = data,
            menu = 0
        ))