/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        "./templates/**/*.html.twig"
    ],
    theme: {
        extend: {
            colors: {
                defaultAccent: '#00b33c',
                pillInc: 'darkseagreen',
                pillOut: 'darksalmon'
            }
        }
    },
    plugins: [
        require('@tailwindcss/forms'),
    ],
}
