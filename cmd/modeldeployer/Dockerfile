# Use the NVIDIA CUDA base image
FROM nvcr.io/nvidia/cuda:12.3.1-base-ubuntu22.04

# Set the working directory
WORKDIR /app

# Install Python and pip
RUN apt-get update && \
    apt-get install -y python3 python3-pip

# Copy the Python files into the container
COPY cmd/modeldeployer/requirements.txt /requirements.txt
COPY cmd/modeldeployer/main.py .

RUN pip install -r /requirements.txt

ENTRYPOINT ["python3", "main.py"]
