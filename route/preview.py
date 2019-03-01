from .tool.func import *

def preview_2(conn, name):
    curs = conn.cursor()

    if acl_check(name) == 1:
        return re_error('/ban')
         
    new_data = re.sub('^\r\n', '', flask.request.form.get('content', None))
    new_data = re.sub('\r\n$', '', new_data)
    
    end_data = render_set(
        title = name,
        data = new_data
    )
    
    if flask.request.args.get('section', None):
        action = '?section=' + flask.request.args.get('section', None)
    else:
        action = ''

    js_data = edit_help_button()
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [name, wiki_set(), custom(), other2([' (' + load_lang('preview') + ')', 0])],
        data =  '<a href="/edit_filter">(' + load_lang('edit_filter_rule') + ')</a>' + js_data[0] + '''
                <form method="post" action="/edit/''' + url_pas(name) + action + '''">
                    ''' + js_data[1] + '''
                    <textarea id="content" rows="25" name="content">''' + html.escape(flask.request.form.get('content', None)) + '''</textarea>
                    <textarea style="display: none;" name="otent">''' + html.escape(flask.request.form.get('otent', None)) + '''</textarea>
                    <hr class=\"main_hr\">
                    <input placeholder="''' + load_lang('why') + '''" name="send" type="text">
                    <hr class=\"main_hr\">
                    ''' + captcha_get() + '''
                    <button id="save" type="submit">''' + load_lang('save') + '''</button>
                    <button id="preview" type="submit" formaction="/preview/''' + url_pas(name) + action + '">' + load_lang('preview') + '''</button>
                </form>
                <hr class=\"main_hr\">
                ''' + end_data,
        menu = [['w/' + url_pas(name), load_lang('return')]]
    ))