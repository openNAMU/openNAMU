import json

print('name : ', end = '')
b = input()

json_data = open('en-US.json', 'rt', encoding='utf-8').read()
a_json = json.loads(json_data)

json_data = open(str(b) + '.json', 'rt', encoding='utf-8').read()
b_json = json.loads(json_data)

for a_in in a_json:
    if not a_in in b_json:
        print(a_in + ' : ', end = '')
        c = input()
        b_json[a_in] = c
        
print(str(b_json).replace(', ', ',\n    ').replace('{', '{\n    ').replace('}', '\n}').replace('\'', '"'))
