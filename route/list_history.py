from .tool.func import *

from .go_api_list_history import api_list_history
from .go_api_list_recent_change import api_list_recent_change

from .recent_change import recent_change_send_render

async def option_lang(lang_in, lang):
    if lang_in == 'user':
        lang_in = lang.get('user_document')
    elif lang_in in lang:
        lang_in = lang[lang_in]
    
    return lang_in

async def list_history(tool = 'history', num = 1, set_type = 'normal', doc_name = 'Test'):
    with get_db_connect() as conn:
        if flask.request.method == 'POST':
            a = flask.request.form.get('a')
            b = flask.request.form.get('b')

            return redirect(conn, f'/diff/{b}/{a}/{doc_name}')
        else:
            if tool == "history":
                data = await api_list_history(num, set_type, doc_name)

                title = doc_name
                sub = '(' + get_lang(conn, 'history') + ') (' + get_lang(conn, set_type) + ')'
                menu = [['w/' + url_pas(doc_name), get_lang(conn, 'return')], ['history_add/' + url_pas(doc_name), get_lang(conn, 'history_add')], ['history_reset/' + url_pas(doc_name), get_lang(conn, 'history_reset')]]
            else:
                data = await api_list_recent_change(num, set_type, 50, '')

                title = get_lang(conn, 'recent_change')
                sub = '(' + get_lang(conn, set_type) + ')'
                menu = [['other', get_lang(conn, 'return')], ['recent_edit_request', get_lang(conn, 'edit_request')]]

            lang = data["language"]
            auth = data["auth"]
            data = data["data"]
            
            data_html = ''

            if tool == "history":
                option_list = ['normal', 'edit', 'move', 'delete', 'revert', 'r1', 'setting']
                for option in option_list:
                    label = await option_lang(option, lang)
                    data_html += f'<a href="/history_page/1/{option}/{doc_name}">({label})</a> '
            else:
                option_list = ['normal', 'edit', 'move', 'delete', 'revert', 'r1', 'edit_request', 'user', 'file', 'category']
                for option in option_list:
                    label = await option_lang(option, lang)
                    data_html += f'<a href="/recent_change/1/{option}">({label})</a> '

            date_heading = ''
            select = ''

            for for_a in range(len(data)):
                if data[for_a][6] != "" and data[for_a][1] == "":
                    if date_heading != '----':
                        data_html += '<h2>----</h2>'
                        date_heading = '----'

                    data_html += await opennamu_make_list('----', '', '', '')
                    continue

                doc_name = url_pas(data[for_a][1])

                left = '<a href="/w/' + doc_name + '">' + html.escape(data[for_a][1]) + '</a> '
                rev = ''

                if data[for_a][6] != "":
                    rev += f'<span style="color: red;">r{data[for_a][0]}</span>'
                else:
                    rev += f'r{data[for_a][0]}'

                select += f'<option value="{data[for_a][0]}">{data[for_a][0]}</option>'

                if int(data[for_a][0]) > 1:
                    before_rev = str(int(data[for_a][0]) - 1)
                    rev = f'<a href="/diff/{before_rev}/{data[for_a][0]}/{doc_name}">{rev}</a>'

                right = f'<span id="opennamu_list_history_{for_a}_over">'
                right += f'<a id="opennamu_list_history_{for_a}" href="javascript:void(0);">'
                right += '<span class="opennamu_svg opennamu_svg_tool">&nbsp;</span></a>'
                right += f'<span class="opennamu_popup_footnote" id="opennamu_list_history_{for_a}_load" style="display: none;"></span>'
                right += '</span> | '
                right += rev + ' | '

                diff_size = data[for_a][5]
                if diff_size == '0':
                    right += f'<span style="color: gray;">{diff_size}</span>'
                elif '+' in diff_size:
                    right += f'<span style="color: green;">{diff_size}</span>'
                else:
                    right += f'<span style="color: red;">{diff_size}</span>'

                right += ' | '
                right += f'{data[for_a][7]} | '

                edit_type = data[for_a][8] if data[for_a][8] != '' else 'edit'
                right += f'{await option_lang(edit_type, lang)} | '

                time_split = data[for_a][2].split(' ')
                if date_heading != time_split[0]:
                    data_html += f'<h2>{time_split[0]}</h2>'
                    date_heading = time_split[0]

                if len(time_split) > 1:
                    right += time_split[1]

                right += f'<span style="display: none;" id="opennamu_history_tool_{for_a}">'

                right += f'<a href="/raw_rev/{data[for_a][0]}/{doc_name}">{lang["raw"]}</a>'
                right += f' | <a href="/revert/{data[for_a][0]}/{doc_name}">{lang["revert"]} (r{data[for_a][0]})</a>'

                if int(data[for_a][0]) > 1:
                    before_rev = str(int(data[for_a][0]) - 1)
                    right += f' | <a href="/revert/{before_rev}/{doc_name}">{lang["revert"]} (r{before_rev})</a>'
                    right += f' | <a href="/diff/{before_rev}/{data[for_a][0]}/{doc_name}">{lang["compare"]}</a>'

                right += f' | <a href="/history/{doc_name}">{lang["history"]}</a>'

                if auth.get("owner") or auth.get("hidel"):
                    right += f' | <a href="/history_hidden/{data[for_a][0]}/{doc_name}">{lang["hide"]}</a>'

                if auth.get("owner"):
                    right += f' | <a href="/history_delete/{data[for_a][0]}/{doc_name}">{lang["history_delete"]}</a>'
                    right += f' | <a href="/history_send/{data[for_a][0]}/{doc_name}">{lang["send_edit"]}</a>'

                right += '</span>'

                bottom = ''
                if data[for_a][4] != "":
                    bottom = recent_change_send_render(html.escape(data[for_a][4]))

                data_html += await opennamu_make_list(left, right, bottom)

                data_html += f'''
                    <script>
                        document.getElementById('opennamu_list_history_{for_a}').addEventListener("click", function() {{
                            opennamu_do_footnote_popover('opennamu_list_history_{for_a}', '', 'opennamu_history_tool_{for_a}', 'open');
                        }});
                        document.addEventListener("click", function() {{
                            opennamu_do_footnote_popover('opennamu_list_history_{for_a}', '', 'opennamu_history_tool_{for_a}', 'close');
                        }});
                    </script>
                '''

            if tool == "history":
                data_html += get_next_page_bottom(conn, f'/history_page/{{}}/{set_type}/{doc_name}', num, data)
                data_html = (
                    '<form method="post">'
                    f'<select name="a">{select}</select> '
                    f'<select name="b">{select}</select> '
                    f'<button type="submit">{lang["compare"]}</button>'
                    '</form>'
                    '<hr class="main_hr"></hr>' + data_html
                )
            else:
                data_html += get_next_page_bottom(conn, f'/recent_change/{{}}/{set_type}', num, data)

            return easy_minify(conn, flask.render_template(skin_check(conn),
                imp = [title, await wiki_set(), await wiki_custom(conn), wiki_css([sub, 0])],
                data = data_html,
                menu = menu
            ))