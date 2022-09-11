import json
import matplotlib.pyplot as plt
import numpy as np


def plt_madd_per_number_of_concurrent_requests():
    fig, ax = plt.subplots()

    f = open('../results/testMAddPerNumberOfParalleledTests.json')
    data = json.load(f)
    bf_pts = [(p["x"], p["y"]) for p in data['bf']]
    data_as_array = np.array(bf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="bf")

    pf_pts = [(p["x"], p["y"]) for p in data['pf']]
    data_as_array = np.array(pf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="pf")

    ax.legend()
    ax.grid(True)
    plt.title('MADD duration per number of concurrent requests')
    plt.xlabel('Number of concurrent requests')
    plt.ylabel('nanoseconds')
    plt.show()


def plt_madd_per_number_of_items():
    fig, ax = plt.subplots()

    f = open('../results/testMAddPerNumberOfItems.json')
    data = json.load(f)
    bf_pts = [(p["x"], p["y"]) for p in data['bf']]
    data_as_array = np.array(bf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="bf")

    pf_pts = [(p["x"], p["y"]) for p in data['pf']]
    data_as_array = np.array(pf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="pf")

    ax.legend()
    ax.grid(True)
    plt.title('MADD duration per number of items')
    plt.xlabel('Number of items')
    plt.ylabel('nanoseconds')
    plt.show()


def plt_exists_per_number_of_concurrent_requests():
    fig, ax = plt.subplots()

    f = open('../results/testExistsPerNumberOfParalleledTests.json')
    data = json.load(f)
    bf_pts = [(p["x"], p["y"]) for p in data['bf']]
    data_as_array = np.array(bf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="bf")

    cf_pts = [(p["x"], p["y"]) for p in data['cf']]
    data_as_array = np.array(cf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="cf")

    pf_pts = [(p["x"], p["y"]) for p in data['pf']]
    data_as_array = np.array(pf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="pf")

    ax.legend()
    ax.grid(True)
    plt.title('Exists duration per number of concurrent requests')
    plt.xlabel('Number of concurrent requests')
    plt.ylabel('nanoseconds')


def plt_mexists_per_number_of_items():
    fig, ax = plt.subplots()

    f = open('../results/testMExistsPerNumberOfItems.json')
    data = json.load(f)
    bf_pts = [(p["x"], p["y"]) for p in data['bf']]
    data_as_array = np.array(bf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="bf")

    cf_pts = [(p["x"], p["y"]) for p in data['cf']]
    data_as_array = np.array(cf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="cf")

    pf_pts = [(p["x"], p["y"]) for p in data['pf']]
    data_as_array = np.array(pf_pts)
    x, y = data_as_array.T
    ax.scatter(x, y, label="pf")

    ax.legend()
    ax.grid(True)
    plt.title('MExists duration per number of items')
    plt.xlabel('Number of items')
    plt.ylabel('nanoseconds')

if __name__ == '__main__':
    plt_exists_per_number_of_concurrent_requests()
    plt_madd_per_number_of_concurrent_requests()
    plt_madd_per_number_of_items()
    plt_mexists_per_number_of_items()
    plt.show()