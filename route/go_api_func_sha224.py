from .tool.func import *

def api_func_sha224(data = 'Test'):
    with get_db_connect() as conn:
        if platform.system() == 'Linux':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.bin"), data], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.bin"), data], stdout = subprocess.PIPE).communicate()[0]
        else:
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.exe"), data], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.exe"), data], stdout = subprocess.PIPE).communicate()[0]

        data = data.decode('utf8')

        return flask.jsonify({ "data" : data })