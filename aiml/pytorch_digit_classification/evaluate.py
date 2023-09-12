import torch
import torchvision.datasets as datasets
import torchvision.transforms as transforms
from torch.utils.data import DataLoader
from simplenet import SimpleNet

transform = transforms.Compose([transforms.ToTensor(), transforms.Normalize((0.1307,), (0.3081,))])
test_dataset = datasets.MNIST(root='data/', train=False, transform=transform, download=True)
test_loader = DataLoader(test_dataset, batch_size=64, shuffle=False)

input_size = 28 * 28
num_classes = 10

model = SimpleNet(input_size, num_classes)
model.load_state_dict(torch.load('simplenet.pth'))

def check_accuracy(loader, model):
    num_correct = 0
    num_samples = 0
    model.eval()

    with torch.no_grad():
        for x, y in loader:
            x = x.reshape(-1, input_size)
            scores = model(x)
            _, predictions = torch.max(scores, 1)
            num_correct += (predictions == y).sum()
            num_samples += predictions.size(0)

    model.train()
    return num_correct / num_samples

print(f"Accuracy on test dataset: {check_accuracy(test_loader, model):.2f}")

