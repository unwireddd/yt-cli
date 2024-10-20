import re

def parse_file(file_name):
    with open(file_name, 'r') as file:
        lines = file.readlines()

    channels = {}
    for line in lines:
        line = line.strip()
        if line:
            match = re.match(r'([A-Z]+)\s+(.*)', line)
            if match:
                key, value = match.groups()
                channels[key] = value

    return channels

def generate_go_map(channels):
    go_map = "var channels = map[string]string{\n"
    for key, value in channels.items():
        go_map += f'    "{key}": "{value}",\n'
    go_map += "}"
    return go_map

def main():
    file_name = 'channels.txt'
    channels = parse_file(file_name)
    go_map = generate_go_map(channels)
    print(go_map)

if __name__ == "__main__":
    main()