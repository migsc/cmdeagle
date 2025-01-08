// Get arguments from environment variables
const name = process.env.ARGS_NAME || 'World';
const age = process.env.ARGS_AGE;

// Get flags from environment variables
const uppercase = process.env.FLAGS_UPPERCASE === 'true';
const lowercase = process.env.FLAGS_LOWERCASE === 'true';
const repeat = parseInt(process.env.FLAGS_REPEAT || '1');

// Construct base greeting
let greeting = `Hello ${name}!`;
if (age) {
    greeting += ` You are ${age} years old.`;
}

// Apply case transformations
if (uppercase) {
    greeting = greeting.toUpperCase();
} else if (lowercase) {
    greeting = greeting.toLowerCase();
}

// Output greeting with repetition
for (let i = 0; i < repeat; i++) {
    console.log(greeting);
}
