from .tool.func import *

def login_oauth_2(conn, platform, func):
    curs = conn.cursor()

    publish_url = load_oauth('publish_url')
    oauth_data = load_oauth(platform)
    api_url = {}
    data = {
        'client_id' : oauth_data['client_id'],
        'client_secret' : oauth_data['client_secret'],
        'redirect_uri' : publish_url + '/oauth/' + platform + '/callback',
        'state' : 'RAMDOMVALUE'
    }

    if platform == 'discord':
        api_url['redirect'] = 'https://discordapp.com/api/oauth2/authorize'
        api_url['token'] = 'https://discordapp.com/api/oauth2/token'
        api_url['profile'] = 'https://discordapp.com/api/users/@me'
    elif platform == 'naver':
        api_url['redirect'] = 'https://nid.naver.com/oauth2.0/authorize'
        api_url['token'] = 'https://nid.naver.com/oauth2.0/token'
        api_url['profile'] = 'https://openapi.naver.com/v1/nid/me'
    elif platform == 'facebook':
        api_url['redirect'] = 'https://www.facebook.com/v3.1/dialog/oauth'
        api_url['token'] = 'https://graph.facebook.com/v3.1/oauth/access_token'
        api_url['profile'] = 'https://graph.facebook.com/me'
    elif platform == 'kakao':
        api_url['redirect'] = 'https://kauth.kakao.com/oauth/authorize'
        api_url['token'] = 'https://kauth.kakao.com/oauth/token'
        api_url['profile'] = 'https://kapi.kakao.com/v2/user/me'

    if func == 'init':
        if oauth_data['client_id'] == '' or oauth_data['client_secret'] == '':
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('error'), wiki_set(), custom(), other2([0, 0])],
                data = load_lang('oauth_disabled'),
                menu = [['user', load_lang('return')]]
            ))
        elif publish_url == 'https://':
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('error'), wiki_set(), custom(), other2([0, 0])],
                data = load_lang('oauth_setting_not_found'),
                menu = [['user', load_lang('return')]]
            ))

        referrer_re = re.compile(r'(?P<host>^(https?):\/\/([^\/]+))\/(?P<refer>[^\/?]+)')
        if flask.request.referrer != None:
            referrer = referrer_re.search(flask.request.referrer)
            if referrer.group('host') != load_oauth('publish_url'):
                return redirect()
            else:
                flask.session['referrer'] = referrer.group('refer')
        else:
            return redirect()

        flask.session['refer'] = flask.request.referrer

        if platform == 'discord':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&response_type=code&scope=identify'.format(
                data['client_id'],
                data['redirect_uri']
            ))
        elif platform == 'naver':
            return redirect(api_url['redirect'] + '?response_type=code&client_id={}&redirect_uri={}&state={}'.format(
                data['client_id'],
                data['redirect_uri'],
                data['state']
            ))
        elif platform == 'facebook':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&state={}'.format(
                data['client_id'],
                data['redirect_uri'],
                data['state']
            ))
        elif platform == 'kakao':
            return redirect(api_url['redirect'] + '?client_id={}&redirect_uri={}&response_type=code'.format(
                data['client_id'],
                data['redirect_uri']
            ))

    elif func == 'callback':
        code = flask.request.args.get('code')
        state = flask.request.args.get('state')

        if code == None:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data = '<h2>ie_wrong_callback</h2>' + load_lang('ie_wrong_callback'),
                menu = [['user', load_lang('return')]]
            ))

        if platform == 'discord':
            data = {
                'client_id'     : data['client_id'],
                'client_secret' : data['client_secret'],
                'grant_type'    : 'authorization_code',
                'redirect_uri'  : data['redirect_uri'],
                'scope'         : 'identify',
                'code'          : code
            }
            headers = {
                'Content-Type': 'application/x-www-form-urlencoded',
                'User-Agent': 'Mozilla/5.0'
            }
            token_exchange = urllib.request.Request(
                api_url['token'],
                data = bytes(urllib.parse.urlencode(data).encode()),
                headers = headers
            )
            token_result = urllib.request.urlopen(token_exchange).read()
            token_json = json.loads(token_result)

            headers = {
                'User-Agent'    : 'Mozilla/5.0',
                'Authorization' : 'Bearer ' + token_json['access_token']
            }
            profile_exchange = urllib.request.Request(
                api_url['profile'],
                headers = headers
            )
            profile_result =  urllib.request.urlopen(profile_exchange).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)
            stand_json = {
                'id'        : profile_result_json['id'],
                'name'      : profile_result_json['username'] + '#' + profile_result_json['discriminator'],
                'picture'   : profile_result_json['avatar']
            }
        elif platform == 'naver':
            token_access = api_url['token'] + '?grant_type=authorization_code&client_id={}&client_secret={}&code={}&state={}'.format(
                data['client_id'],
                data['client_secret'],
                code,
                state
            )
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            headers = {
                'Authorization': 'Bearer {}'.format(token_result_json['access_token'])
            }

            profile_access = urllib.request.Request(api_url['profile'], headers = headers)
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {
                'id'        : profile_result_json['response']['id'],
                'name'      : profile_result_json['response']['name'],
                'picture'   : profile_result_json['response']['profile_image']
            }
        elif platform == 'facebook':
            token_access = api_url['token'] + '?client_id={}&redirect_uri={}&client_secret={}&code={}'.format(
                data['client_id'],
                data['redirect_uri'],
                data['client_secret'],
                code
            )
            token_result = urllib.request.urlopen(token_access).read().decode('utf-8')
            token_result_json = json.loads(token_result)

            profile_access = api_url['profile'] + '?fields=id,name,picture&access_token={}'.format(token_result_json['access_token'])
            profile_result = urllib.request.urlopen(profile_access).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)

            stand_json = {
                'id': profile_result_json['id'],
                'name': profile_result_json['name'],
                'picture': profile_result_json['picture']['data']['url']
            }
        elif platform == 'kakao':
            data = {
                'client_id'     : data['client_id'],
                'client_secret' : data['client_secret'],
                'grant_type'    : 'authorization_code',
                'redirect_uri'  : data['redirect_uri'],
                'code'          : code
            }
            headers = {
                'Content-Type': 'application/x-www-form-urlencoded',
                'User-Agent': 'Mozilla/5.0'
            }
            token_exchange = urllib.request.Request(
                api_url['token'],
                data = bytes(urllib.parse.urlencode(data).encode()),
                headers = headers
            )
            token_result = urllib.request.urlopen(token_exchange).read()
            token_json = json.loads(token_result)

            headers = {
                'User-Agent'    : 'Mozilla/5.0',
                'Authorization' : 'Bearer ' + token_json['access_token']
            }
            profile_exchange = urllib.request.Request(
                api_url['profile'],
                headers = headers
            )
            profile_result =  urllib.request.urlopen(profile_exchange).read().decode('utf-8')
            profile_result_json = json.loads(profile_result)
            stand_json = {
                'id'        : profile_result_json['id'],
                'name'      : profile_result_json['properties']['nickname'],
                'picture'   : profile_result_json['properties']['profile_image']
            }

        if flask.session['referrer'][0:6] == 'change':
            curs.execute(db_change('select * from oauth_conn where wiki_id = ? and provider = ?'), [flask.session['id'], platform])
            oauth_result = curs.fetchall()
            if len(oauth_result) == 0:
                curs.execute(db_change('insert into oauth_conn (provider, wiki_id, sns_id, name, picture) values(?, ?, ?, ?, ?)'), [
                    platform,
                    flask.session['id'],
                    stand_json['id'],
                    stand_json['name'],
                    stand_json['picture']
                ])
            else:
                curs.execute(db_change('update oauth_conn set name = ? picture = ? where wiki_id = ?'), [
                    stand_json['name'],
                    stand_json['picture'],
                    flask.session['id']
                ])

            conn.commit()
        elif flask.session['referrer'][0:5] == 'login':
            curs.execute(db_change('select * from oauth_conn where provider = ? and sns_id = ?'), [platform, stand_json['id']])
            curs_result = curs.fetchall()
            if len(curs_result) == 0:
                return re_error('/error/2')
            else:
                flask.session['state'] = 1
                flask.session['id'] = curs_result[0][2]

        return redirect(flask.session['refer'])