from .tool.func import *

class api_w_render_include:
    def __init__(self, data_option):
        self.include_change_list = data_option

    def __call__(self, match):
        match_org = match.group(0)
        match = match.groups()

        if len(match) < 3:
            match = list(match) + ['']

        if match[2] == '\\':
            return match_org
        else:
            slash_add = ''
            if match[0]:
                if len(match[0]) % 2 == 1:
                    slash_add = '\\' * (len(match[0]) - 1)
                else:
                    slash_add = match[0]

            if match[1] in self.include_change_list:
                return slash_add + self.include_change_list[match[1]]
            else:
                return slash_add + match[2]

def api_w_render(name = '', tool = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            name = flask.request.form.get('name', '')
            data_org = flask.request.form.get('data', '')
            
            data_option = flask.request.form.get('option', '')
            if data_option != '':
                data_option = json.loads(data_option)

                data_option_func = api_w_render_include(data_option)

                # parameter replace
                data_org = re.sub(r'(\\+)?@([ㄱ-힣a-zA-Z0-9]+)=((?:\\@|[^@\n])+)@', data_option_func, data_org)
                data_org = re.sub(r'(\\+)?@([ㄱ-힣a-zA-Z0-9]+)@', data_option_func, data_org)

                # remove end br
                data_org = re.sub('^\n+', '', data_org)

            if tool == '':
                data_pas = render_set(
                    doc_name = name, 
                    doc_data = data_org, 
                    data_type = 'api_view'
                )
            elif tool == 'include':
                data_pas = render_set(
                    doc_name = name, 
                    doc_data = data_org, 
                    data_type = 'api_include'
                )
            else:
                data_pas = render_set(
                    doc_name = name,
                    doc_data = data_org, 
                    data_type = 'api_thread'
                )

            return flask.jsonify({
                "data" : data_pas[0], 
                "js_data" : data_pas[1]
            })
        else:
            return ''