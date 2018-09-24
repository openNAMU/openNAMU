from . import tool

import datetime
import html
import re

def html_only(conn, data, title, main_num):
	curs = conn.cursor()
	
	backlink = []
	plus_data = ''

	while 1:
		in_data = re.search('<a(?: href="/w/((?:(?!\").)+)")?>((?:(?!<\/a>).)+)<\/a>', data)
		if in_data:
			in_data = in_data.groups()

			if in_data[0]:
				main_link = in_data[0]
				sub_link = in_data[1]
			else:
				main_link = in_data[1]
				sub_link = in_data[1]

			curs.execute("select title from data where title = ?", [main_link])
			if not curs.fetchall():
				link_id = 'id="not_thing"'
			
				backlink += [[title, main_link, 'no']]
			else:
				link_id = 'id=""'

			backlink += [[title, main_link, '']]

			data = re.sub('<a(?: href="/((?:(?!\").)+)")?>((?:(?!<\/a>).)+)<\/a>', '<a ' + link_id + ' href="/w/' + main_link + '">' + sub_link + '</a>', data, 1)
		else:
			break

	data = re.sub('<test_a', '', data)

	while 1:
		in_data = re.search('<a((?:(?!>).)+)>((?:(?!<\/a>).)+)<\/a>', data)
		if in_data:
			in_data = in_data.groups()

			a_data = re.sub('href="((?:(?!").)+)"', '', in_data[0])
			a_data = re.sub('id="((?:(?!").)+)"', '', a_data)

			if re.search('=', a_data):
				data = re.sub('<a((?:(?!>).)+)>((?:(?!<\/a>).)+)<\/a>', '', data, 1)
			else:
				data = re.sub('<a((?:(?!>).)+)>((?:(?!<\/a>).)+)<\/a>', '<test_a ' + in_data[0] + '>' + in_data[1] + '</a>', data, 1)
		else:
			break

	data = re.sub('<test_a', '<a', data)

	ok_list = [
		'h1',
		'h2',
		'h3',
		'h4',
		'h5',
		'h6',
		'li',
		'a'
	]

	data = tool.xss_protect(curs, data, ok_list)
	
	return [data, plus_data, backlink]
