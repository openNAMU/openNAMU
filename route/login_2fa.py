from .tool.func import *
import hmac
import hashlib
import os
import time
import math
import base64
import urllib.parse
import random

import segno

def login_2fa_2(conn):
    curs = conn.cursor()

    if not (flask.session and 'b_id' in flask.session):
        return redirect('/user')

    ip = ip_check()
    if ip_or_user(ip) == 0:
        return redirect('/user')

    if ban_check(None, 'login') == 1:
        return re_error('/ban')

    if flask.request.method == 'POST':
        if captcha_post(flask.request.form.get('g-recaptcha-response', flask.request.form.get('g-recaptcha', ''))) == 1:
            return re_error('/error/13')
        else:
            captcha_post('', 0)

        user_agent = flask.request.headers.get('User-Agent', '')
        user_id = flask.session['b_id']

        curs.execute(db_change('select data from user_set where name = "2fa_pw" and id = ?'), [user_id])
        user_1 = curs.fetchall()
        if user_1:
            curs.execute(db_change('select data from user_set where name = "2fa_pw_encode" and id = ?'), [user_id])
            user_1 = user_1[0][0]
            user_2 = curs.fetchall()[0][0]
            pw_check_d = 0
            if (user_2 == "totp"):
                code = dev_2fa_totp(user_1)
                time.sleep(random.random() * 0.5)
                if (code != flask.request.form.get('pw', '')):
                    pw_check_d = 1
            else:
                pw_check_d = pw_check(
                    flask.request.form.get('pw', ''),
                    user_1,
                    user_2,
                    user_id
                )
            if pw_check_d != 1:
                return re_error('/error/10')

        flask.session['id'] = user_id

        ua_plus(user_id, ip, user_agent, get_time())
        conn.commit()

        flask.session.pop('b_id', None)

        return redirect('/user')
    else:
        return easy_minify(flask.render_template(skin_check(),
            imp = [load_lang('login'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <input placeholder="''' + load_lang('2fa_password') + '''" name="pw" type="password">
                        <hr class=\"main_hr\">
                        ''' + captcha_get() + '''
                        <button type="submit">''' + load_lang('login') + '''</button>
                        ''' + http_warring() + '''
                    </form>
                    ''',
            menu = [['user', load_lang('return')]]
        ))

def dev_2fa_totp(key):
    key = base64.b32decode(key)
    counter = math.floor(time.time()/30)
    b = counter.to_bytes(8, "big")
    #print(list(b))
    mac = hmac.new(key, digestmod=hashlib.sha1)
    mac.update(b)
    hmacResult = mac.digest()
    offset = int(hmacResult[19] & 0xf)
    bincode = int(
		(int(hmacResult[offset]&0x7f))<<24 |
		(int(hmacResult[offset+1]&0xff))<<16 |
		(int(hmacResult[offset+2]&0xff))<<8 |
		(int(hmacResult[offset+3] & 0xff)),
	)
    code = str(bincode%1000000)
    return (6-len(code)) * "0" + code

def dev_2fa_new_key():
    secret = os.urandom(20)
    key = base64.b32encode(secret)
    return key.decode("utf-8")

def dev_2fa_qrcode(provider,accountname,key):
    urlstr = urllib.parse.urlencode({"secret":key, "issuer":provider, "algorithm": "SHA1", "digits": "6", "period": "30"})
    #print(urlstr)
    qrstr = "otpauth://totp/{}:{}?{}".format(urllib.parse.quote(provider), urllib.parse.quote(accountname), urlstr)
    qr = segno.make_qr(qrstr, error="H")
    return qr.svg_data_uri(scale=2)

# if __name__ == "__main__":
#     key = dev_2fa_new_key()
#     print(key)
#     print(dev_2fa_qrcode("example.com 위키", "admin", key))
#     print(dev_2fa_totp(key))
