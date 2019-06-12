from .tool.func import *

def setting_oauth_2(conn):
    curs = conn.cursor()

    if admin_check(None, 'oauth setting') != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        try:
            facebook_client_id = flask.request.form['facebook_client_id']
            facebook_client_secret = flask.request.form['facebook_client_secret']
            
            naver_client_id = flask.request.form['naver_client_id']
            naver_client_secret = flask.request.form['naver_client_secret']
        except:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data = '<h2>ie_no_data_required</h2>' + load_lang('ie_no_data_required'),
                menu = [['other', load_lang('return')]]
            ))

        with open(app_var['path_oauth_setting'], 'r', encoding='utf-8') as f:
            legacy = json.loads(f.read())

        with open(app_var['path_oauth_setting'], 'w', encoding='utf-8') as f:
            f.write('''
                {
                    "_README" : {
                        "en" : "''' + legacy['_README']['en'] + '''",
                        "ko" : "''' + legacy['_README']['ko'] + '''",
                        "support" : ''' + str(legacy['_README']['support']).replace("'", '"') + '''
                    },
                    "publish_url" : "''' + legacy['publish_url'] + '''",
                    "facebook" : {
                        "client_id" : "''' + facebook_client_id + '''",
                        "client_secret" : "''' + facebook_client_secret + '''"
                    },
                    "naver" : {
                        "client_id" : "''' + naver_client_id + '''",
                        "client_secret" : "''' + naver_client_secret + '''"
                    }
                }
            ''')
        
        return flask.redirect('/oauth_setting')

    oauth_supported = load_oauth('_README')['support']

    body_content = ''
    body_content += '''
        <script>
            function check_value (target) {
                target_box = document.getElementById(target.id + "_box");
                if (target.value !== "") {
                    target_box.checked = true;
                } else {
                    target_box.checked = false;
                } 
            }
        </script>
    '''

    init_js = ''
    body_content += '<form method="post">'

    for i in range(len(oauth_supported)):
        oauth_data = load_oauth(oauth_supported[i])
        for j in range(2):
            if j == 0:
                load_target = 'id'
            elif j == 1:
                load_target = 'secret'

            init_js += 'check_value(document.getElementById("{}_client_{}"));'.format(oauth_supported[i], load_target)

            body_content += '''
                <input id="{}_client_{}_box" type="checkbox" disabled>
                <input placeholder="{}_client_{}" id="{}_client_{}" name="{}_client_{}" value="{}" type="text" onChange="check_value(this)" style="width: 80%;">
                <hr>
            '''.format(
                oauth_supported[i],
                load_target,
                oauth_supported[i], 
                load_target, 
                oauth_supported[i], 
                load_target, 
                oauth_supported[i], 
                load_target, 
                oauth_data['client_{}'.format(load_target)]
            )
    
    body_content += '<button id="save" type="submit">' + load_lang('save') + '</button></form>'
    body_content += '<script>' + init_js + '</script>'
    
    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('oauth_setting'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('return')]]
    ))