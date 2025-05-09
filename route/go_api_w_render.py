from .tool.func import *

from .go_api_w_raw import api_w_raw

async def api_w_render(name = '', tool = '', request_method = '', request_data = {}):
    with get_db_connect() as conn:
        curs = conn.cursor()

        flask_data = flask_data_or_variable(flask.request.form, request_data)
        request_method = flask.request.method if request_method == '' else request_method

        if request_method == 'POST':
            name = flask_data.get('name', '')
            data_org = flask_data.get('data', '')
            data_option = flask_data.get('option', '')

            markup = ''
            if tool in ('', 'from', 'include'):
                curs.execute(db_change("select set_data from data_set where doc_name = ? and set_name = 'document_markup'"), [name])
                db_data = curs.fetchall()
                if db_data and db_data[0][0] != '' and db_data[0][0] != 'normal':
                    markup = db_data[0][0]

                if markup == '':
                    curs.execute(db_change('select data from other where name = "markup"'))
                    db_data = curs.fetchall()
                    markup = db_data[0][0] if db_data else 'namumark'

            data_type = ''
            if tool == '':
                data_type = 'api_view'
            elif tool == 'from':
                data_type = 'api_from'
            elif tool == 'include':
                data_type = 'api_include'
            elif tool == 'backlink':
                data_type = 'backlink'
            else:
                data_type = 'api_thread'

            if markup in ('', 'namumark', 'namumark_beta') and data_option != '':
                data_option = json_loads(data_option)

                # remove end br
                data_org = re.sub('^\n+', '', data_org)

            if markup in ('', 'namumark'):
                data_pas = render_set(conn, 
                    doc_name = name, 
                    doc_data = data_org, 
                    data_type = data_type,
                    parameter = data_option
                )

                return {
                    "data" : data_pas[0], 
                    "js_data" : data_pas[1]
                }
            else:
                other_set = {}
                other_set["doc_name"] = name
                other_set["render_type"] = data_type
                other_set["data"] = data_org

                return await python_to_golang(sys._getframe().f_code.co_name, other_set)
        else:
            return {}

async def api_w_render_exter(name = '', tool = '', request_method = '', request_data = {}):
    return flask.jsonify(await api_w_render(name, tool, request_method, request_data))