import torch
import torch.nn as nn
import torch.optim as optim
import torchvision.datasets as datasets
import torchvision.transforms as transforms
from torch.utils.data import DataLoader
from simplenet import SimpleNet

transform = transforms.Compose([transforms.ToTensor(), transforms.Normalize((0.1307,), (0.30811,))])

train_dataset = datasets.MNIST(root='data/', train=True, transform=transform, download=True)
train_dataset = datasets.MNIST(root='data/', train=True, transform=transform, download=True)
train_loader = DataLoader(train_dataset, batch_size=64, shuffle=True)

input_size = 28 * 28
num_classes = 10
learning_rate = 0.001
num_epochs = 10

model = SimpleNet(input_size, num_classes)
criterion = nn.CrossEntropyLoss()
optimizer = optim.Adam(model.parameters(), lr=learning_rate)

for epoch in range(num_epochs):
    for batch_idx, (data, targets) in enumerate(train_loader):
        data = data.reshape(-1, input_size)
        scores = model(data)
        loss = criterion(scores, targets)

        optimizer.zero_grad()
        loss.backward()
        optimizer.step()

                       
