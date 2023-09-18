from .tool.func import *

def view_random(db_set):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if platform.system() == 'Linux':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.bin"), db_set], stdout=subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.bin"), db_set], stdout=subprocess.PIPE).communicate()[0]
        else:
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.exe"), db_set], stdout=subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.exe"), db_set], stdout=subprocess.PIPE).communicate()[0]

        return redirect('/w/' + url_pas(data.decode('utf8')))