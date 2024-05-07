# Use the official Python image as the base image
FROM python:3.10.12

# Set the working directory inside the container
WORKDIR /app

# Copy the requirements file
COPY requirements.txt .

# Install Python dependencies
RUN pip install -r requirements.txt

ENV CI_TOKEN=$CI_TOKEN
# Copy the Flask app code to the container
COPY . .

# Expose the Flask app port

# Command to run the Flask app
CMD ["python", "CD_service.py"]