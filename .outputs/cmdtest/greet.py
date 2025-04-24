import os

# Get arguments from environment variables
name = os.environ.get('ARGS_NAME', 'World')
age = os.environ.get('ARGS_AGE')

# Get flags from environment variables
uppercase = os.environ.get('FLAGS_UPPERCASE', 'false').lower() == 'true'
lowercase = os.environ.get('FLAGS_LOWERCASE', 'false').lower() == 'true'
repeat = int(os.environ.get('FLAGS_REPEAT', '1'))

# Construct base greeting
greeting = f"Hello {name}!"
if age:
    greeting += f" You are {age} years old."

# Apply case transformations
if uppercase:
    greeting = greeting.upper()
elif lowercase:
    greeting = greeting.lower()

# Output greeting with repetition
for _ in range(repeat):
    print(greeting)
