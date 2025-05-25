#!/usr/bin/env python3
import matplotlib.pyplot as plt
import numpy as np
import random

def plot(means):
    plt.hist(means, bins=30, density=True, alpha=0.6, color='g')
    plt.title("Standard Normal Distribution (mean=0, variance=1)")
    plt.xlabel("Value")
    plt.ylabel("Probability Density")
    plt.show()

class Bandit:
    def __init__(self, k=10):
        self.q = np.random.normal(loc=0, scale=1, size=k)

    def bandit(self, action):
        return np.random.normal(self.q[action], 1)

class egreedy:
    def __init__(self, epsilon, k):
        self.epsilon = epsilon
        self.Q = np.zeros(k)
        self.N = np.zeros(k)
        self.bandit = Bandit(k)
        self.step_number = 0
        self.k = k

    @staticmethod
    def desc(epsilon, k):
        return f"epsilon greedy - epsilon {epsilon}, k {k}"

    def chooseQ(self):
        if np.random.rand() < self.epsilon:
            # epsilon case choose from all actions
            return np.random.randint(self.k)
        else:
            return np.argmax(self.Q)

    def step(self):
        action = self.chooseQ()
        reward = self.bandit.bandit(action)
        self.N[action] += 1
        self.Q[action] = self.Q[action] + ((reward - self.Q[action])/self.N[action])
        return action, reward

class reducing_egreedy:
    def __init__(self, k):
        self.epsilon = 1
        self.Q = np.zeros(k)
        self.N = np.zeros(k)
        self.bandit = Bandit(k)
        self.step_number = 0
        self.k = k
        self.nvalue = 0

    @staticmethod
    def desc(k):
        return f"reducing epsilon greedy - k {k}"

    def chooseQ(self):
        if np.random.rand() < self.epsilon:
            # epsilon case choose from all actions
            return np.random.randint(self.k)
        else:
            return np.argmax(self.Q)

    def step(self):
        self.step_number += 1
        if self.step_number >= (self.nvalue + 1) * 50:
            self.nvalue+=1
            self.epsilon = (1/(1 + self.nvalue))
        action = self.chooseQ()
        reward = self.bandit.bandit(action)
        self.N[action] += 1
        self.Q[action] = self.Q[action] + ((reward - self.Q[action])/self.N[action])
        return action, reward

class ucb:
    def __init__(self, c, k):
        self.c = c
        self.Q = np.zeros(k)
        self.N = np.zeros(k)
        self.bandit = Bandit(k)
        self.step_number = 0
        self.k = k

    @staticmethod
    def desc(c, k):
        return f"ucb - c {c}, k {k}"

    def chooseQ(self):
        zero_N = np.where(self.N == 0)[0]
        if len(zero_N) > 0:
            return zero_N[0]  # Prefer unexplored actions

        confidence_bounds = self.c * np.sqrt(np.emath.logn(np.e, self.step_number) / self.N)
        values = self.Q + confidence_bounds
        return np.argmax(values)

    def step(self):
        self.step_number += 1
        action = self.chooseQ()
        reward = self.bandit.bandit(action)
        self.N[action] += 1
        self.Q[action] = self.Q[action] + ((reward - self.Q[action])/self.N[action])
        return action, reward

def run_bandit_problem(algo, args, timesteps, runs):
    optimal_action_counts = np.zeros(timesteps)
    average_rewards = np.zeros(timesteps)

    for run in range(runs):
        print("\r",run,end="")
        inst = algo(*args)
        optimal_action = np.argmax(inst.bandit.q)
        rewards = np.zeros(timesteps)

        for t in range(timesteps):
            action, reward = inst.step()
            if action == optimal_action:
                optimal_action_counts[t] += 1
            rewards[t] = reward
        average_rewards += rewards

    average_rewards /= runs
    optimal_action_counts /= runs
    return average_rewards, optimal_action_counts

algos = [ egreedy, egreedy, reducing_egreedy, ucb, ucb]
args = [[ 0.1, 10 ], [0.2, 10], [10], [1, 10], [10, 10]]
figs = []
for idx, algo in enumerate(algos):
    print(algo.desc(*args[idx]))
    avg_rewards, opt_actions = run_bandit_problem(algo, args[idx], 1000, 2000)
    fig, axs = plt.subplots(2, 1)

    # adjusting for 1 indexing.
    avg_rewards = np.insert(avg_rewards, 0, 0)
    opt_actions = np.insert(opt_actions, 0, 0)

    axs[0].plot(avg_rewards)
    axs[0].set_xlabel('Steps')
    axs[0].set_ylabel('Average Reward')
    axs[0].set_ylim(0, 2)
    axs[1].plot(opt_actions * 100)
    axs[1].set_xlabel('Steps')
    axs[1].set_ylabel('% Optimal Action')
    axs[1].set_ylim(0, 100)
    fig.suptitle(algo.desc(*args[idx]))
    figs.append(fig)
    print()
print()
plt.show()
